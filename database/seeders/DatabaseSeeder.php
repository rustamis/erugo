<?php

namespace Database\Seeders;

use App\Models\User;
// use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class DatabaseSeeder extends Seeder
{
    /**
     * Seed the application's database.
     */
    public function run(): void
    {
        $this->call([
            SettingsSeeder::class,
        ]);

        // load themes.sql from the root of the project
        $sql = file_get_contents(base_path('themes.sql'));
        //run the query, ignore errors
        DB::unprepared($sql);
    }
}
