<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\Auth;
use App\Models\Share;
use App\Models\File;
use Illuminate\Support\Str;
use App\Jobs\CreateShareZip;
use Carbon\Carbon;
use App\Haikunator;
use App\Models\Setting;
use App\Models\User;
use App\Models\Download;
class SharesController extends Controller
{
    public function create(Request $request)
    {
        $validator = Validator::make($request->all(), [
            'name' => ['string', 'max:255'],
            'description' => ['max:255'],
            'expires_at' => ['date'],
            'files' => ['required', 'array']
        ]);

        if ($validator->fails()) {
            return response()->json([
                'status' => 'error',
                'message' => 'Validation failed',
                'data' => [
                    'errors' => $validator->errors()
                ]
            ], 422);
        }

        $user = Auth::user();

        if (!$user) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }
        $longId = $this->generateLongId();
        $files = $request->file('files');
        $totalFileSize = 0;
        foreach ($files as $file) {
            $totalFileSize += $file->getSize();
        }

        $sharePath = storage_path('app/shares/' . $user->id . '/' . $longId);

        $shareData = [
            'name' => $request->name,
            'description' => $request->description,
            'expires_at' => Carbon::now()->addDays(7),
            'user_id' => $user->id,
            'path' => $sharePath,
            'long_id' => $longId,
            'size' => $totalFileSize,
            'file_count' => count($files)
        ];
        $share = Share::create($shareData);
        foreach ($files as $file) {
            $fileData = [
                'share_id' => $share->id,
                'name' => $file->getClientOriginalName(),
                'type' => $file->getMimeType(),
                'size' => $file->getSize()
            ];
            $file = File::create($fileData);
            $file->share_id = $share->id;
            $file->save();
        }

        if (!file_exists($sharePath)) {
            mkdir($sharePath, 0777, true);
        }
        $files = $request->file('files');
        foreach ($files as $file) {
            $file->move($sharePath, $file->getClientOriginalName());
        }
        $share->status = 'pending';
        $share->save();

        //dispatch the job to create the zip file
        CreateShareZip::dispatch($share);

        return response()->json([
            'status' => 'success',
            'message' => 'Files uploaded successfully',
            'data' => [
                'share' => $share
            ]
        ]);
    }

    public function read($shareId)
    {
        $share = Share::where('long_id', $shareId)->with('files')->first();
        if (!$share) {
            return response()->json([
                'status' => 'error',
                'message' => 'Share not found'
            ], 404);
        }

        if ($share->expires_at < Carbon::now()) {
            return response()->json([
                'status' => 'error',
                'message' => 'Share expired'
            ], 410);
        }

        if ($share->download_limit != null && $share->download_count >= $share->download_limit) {
            return response()->json([
                'status' => 'error',
                'message' => 'Download limit reached'
            ], 410);
        }


        return response()->json([
            'status' => 'success',
            'message' => 'Share found',
            'data' => [
                'share' => $share
            ]
        ]);
    }

    public function download($shareId)
    {

        $share = Share::where('long_id', $shareId)->with('files')->first();
        if (!$share) {
            return redirect()->to('/shares/' . $shareId);
        }

        if ($share->expires_at < Carbon::now()) {
            return redirect()->to('/shares/' . $shareId);
        }

        if ($share->download_limit != null && $share->download_count >= $share->download_limit) {
            return redirect()->to('/shares/' . $shareId);
        }

        //if there is only one file, download it directly
        if ($share->file_count == 1) {
            if (file_exists($share->path . '/' . $share->files[0]->name)) {
                $share->download_count++;
                $share->save();

                $this->createDownloadRecord($share);

                return response()->download($share->path . '/' . $share->files[0]->name);
            } else {
                return redirect()->to('/shares/' . $shareId);
            }
        }

        //otherise let's check the status: pending, ready, or failed
        if ($share->status == 'pending') {
            return view('shares.pending', [
                'share' => $share,
                'settings' => $this->getSettings()
            ]);
        }

        //if the share is ready, download the zip file
        if ($share->status == 'ready') {
            $filename = $share->path . '.zip';
            //does the file exist?
            if (file_exists($filename)) {
                //we're as sure as we can be that this download is going to happen, so let's mark the share as downloaded +1 time!
                $share->download_count++;
                $share->save();

                $this->createDownloadRecord($share);
                
                return response()->download($filename);
            } else {
                //something went wrong, show the failed view
                return view('shares.failed', [
                    'share' => $share,
                    'settings' => $this->getSettings()
                ]);
            }
        }

        //if the share is failed, show the failed view
        if ($share->status == 'failed') {
            return view('shares.failed', [
                'share' => $share,
                'settings' => $this->getSettings()
            ]);
        }

        //if we got here we have no idea what to do so let's show the failed view
        return view('shares.failed', [
            'share' => $share,
            'settings' => $this->getSettings()
        ]);
    }

    public function myShares()
    {
        $user = Auth::user();

        if (!$user) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }

        $shares = Share::where('user_id', $user->id)->orderBy('created_at', 'desc')->with('files')->get();
        return response()->json([
            'status' => 'success',
            'message' => 'My shares',
            'data' => [
                'shares' => $shares
            ]
        ]);
    }

    public function expire($shareId)
    {
        $user = Auth::user();
        if (!$user) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }
        $share = Share::where('id', $shareId)->first();
        if (!$share) {
            return response()->json([
                'status' => 'error',
                'message' => 'Share not found'
            ], 404);
        }
        if($share->user_id != $user->id) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }
        $share->expires_at = Carbon::now();
        $share->save();

        return response()->json([
            'status' => 'success',
            'message' => 'Share expired',
            'data' => [
                'share' => $share
            ]
        ]);
    }

    public function extend($shareId)
    {

        $user = Auth::user();
        if (!$user) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }
        $share = Share::where('id', $shareId)->first();
        if (!$share) {
            return response()->json([
                'status' => 'error',
                'message' => 'Share not found'
            ], 404);
        }
        if($share->user_id != $user->id) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }
        $share->expires_at = Carbon::now()->addDays(7);
        $share->save();
        return response()->json([
            'status' => 'success',
            'message' => 'Share extended',
            'data' => [
                'share' => $share
            ]
        ]);
    }

    public function setDownloadLimit($shareId, Request $request)
    {
        $user = Auth::user();
        if (!$user) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }
        $share = Share::where('id', $shareId)->first();
        if (!$share) {
            return response()->json([
                'status' => 'error',
                'message' => 'Share not found'
            ], 404);
        }
        if($share->user_id != $user->id) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }
        if($request->amount == -1) {
            $share->download_limit = null;
        } else {
            $share->download_limit = $request->amount;
        }
        $share->save();
        return response()->json([
            'status' => 'success',
            'message' => 'Download limit set',
            'data' => [
                'share' => $share
            ]
        ]);
    }


    private function generateLongId()
    {
        $maxAttempts = 10;
        $attempts = 0;
        $id = Haikunator::haikunate() . '-' . Haikunator::haikunate();
        while (Share::where('long_id', $id)->exists() && $attempts < $maxAttempts) {
            $id = Haikunator::haikunate() . '-' . Haikunator::haikunate();
            $attempts++;
        }
        if ($attempts >= $maxAttempts) {
            throw new \Exception('Unable to generate unique long_id after ' . $maxAttempts . ' attempts');
        }
        return $id;
    }

    private function getSettings()
    {
        $settings = Setting::whereLike('group', 'ui%')->get();
        $indexedSettings = [];
        foreach ($settings as $setting) {
            $indexedSettings[$setting->key] = $setting->value;
        }

        //have we any users in the database?
        $userCount = User::count();
        $indexedSettings['setup_needed'] = $userCount > 0 ? 'false' : 'true';

        //grab the app url from env
        $appURL = env('APP_URL');
        $indexedSettings['api_url'] = $appURL;

        return $indexedSettings;
    }

    private function createDownloadRecord(Share $share) {
        $ipAddress = request()->ip();
        $userAgent = request()->userAgent();
        $download = Download::create([
            'share_id' => $share->id,
            'ip_address' => $ipAddress,
            'user_agent' => $userAgent
        ]);
        $download->save();
        return $download;
    }
}
