<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;
use App\Services\SettingsService;

class maxRequestSize
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {

        $maxUploadSize = app(SettingsService::class)->getMaxUploadSize();
        $maxUploadSizeValue = app(SettingsService::class)->get('max_share_size');
        $maxUploadSizeUnit = app(SettingsService::class)->get('max_share_size_unit');
    

        if ($request->hasHeader('Content-Length') && 
            $request->header('Content-Length') > $maxUploadSize) {
            
            if ($request->expectsJson()) {
                return response()->json([
                    'error' => 'Total upload size exceeds maximum allowed size of ' . $maxUploadSizeValue . $maxUploadSizeUnit
                ], 413);
            }

            return redirect()->back()
                ->withErrors(['files' => 'Total upload size exceeds maximum allowed size of ' . $maxUploadSizeValue . $maxUploadSizeUnit]);
        }

        return $next($request);
    }
}
