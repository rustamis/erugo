<?php

namespace App\Providers;

use Illuminate\Support\ServiceProvider;
use App\Services\SettingsService;
use Illuminate\Support\Facades\View;
use Illuminate\Support\Facades\Schema;
class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        //
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        
        //only run this if the database has been set up
        if (Schema::hasTable('settings')) {
            $settingsService = new SettingsService();
            $settings = $settingsService->getGlobalViewData();
            View::share('settings', $settings);
        }


        View::prependLocation(storage_path('templates'));
    }
}
