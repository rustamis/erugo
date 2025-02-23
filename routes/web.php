<?php

use Illuminate\Support\Facades\Route;
use App\Models\Setting;
use Illuminate\Support\Facades\Storage;
use App\Models\User;

Route::get('/', function () {
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

    return view('app', ['settings' => $indexedSettings]);
});

Route::get('/test', function () {
    return response()->json([
        'app_url' => env('APP_URL'),
    ]);
});

Route::get('/shares/{share}', function () {
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

    return view('app', ['settings' => $indexedSettings]);
});


Route::get('/logo', function () {
    //grab the logo file data from settings
    $setting = Setting::where('key', 'logo')->first();
    $logo = Storage::disk('public')->get($setting->value);
    // return $setting;
    return response($logo)->header('Content-Type', 'image/png');
});


Route::get('/info',function() {

    return response()->json([
        'app_url' => env('APP_URL'),
    ]);

});