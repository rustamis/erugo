<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class SettingsSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('settings')->insertOrIgnore([

            [
                'key' => 'application_name',
                'value' => 'Erugo File Sharing',
                'previous_value' => 'Erugo File Sharing',
                'group' => 'ui.strings',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'css_primary_color',
                'value' => '#589db6',
                'previous_value' => '#589db6',
                'group' => 'ui.css',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'css_secondary_color',
                'value' => '#01021c',
                'previous_value' => '#01021c',
                'group' => 'ui.css',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'css_accent_color',
                'value' => '#63a8bc',
                'previous_value' => '#63a8bc',
                'group' => 'ui.css',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'css_accent_color_light',
                'value' => '#d0e1d5',
                'previous_value' => '#d0e1d5',
                'group' => 'ui.css',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'logo',
                'value' => 'logo.png',
                'previous_value' => null,
                'group' => 'ui.logo',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'logo_width',
                'value' => '100',
                'previous_value' => '100',
                'group' => 'ui.logo',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'use_my_backgrounds',
                'value' => 'false',
                'previous_value' => 'false',
                'group' => 'ui',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'show_powered_by',
                'value' => 'true',
                'previous_value' => 'true',
                'group' => 'ui',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'application_url',
                'value' => 'http://localhost:9199',
                'previous_value' => null,
                'group' => 'system',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'login_message',
                'value' => 'Login to your account to get started.',
                'previous_value' => null,
                'group' => 'ui.strings',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'max_expiry_time',
                'value' => '10',
                'previous_value' => null,
                'group' => 'system.shares',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'max_share_size',
                'value' => '2',
                'previous_value' => null,
                'group' => 'system.shares',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'max_share_size_unit',
                'value' => 'GB',
                'previous_value' => null,
                'group' => 'system.shares',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'clean_files_after_days',
                'value' => '30',
                'previous_value' => null,
                'group' => 'system.shares',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'smtp_host',
                'value' => null,
                'previous_value' => null,
                'group' => 'system.smtp',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'smtp_port',
                'value' => '587',
                'previous_value' => null,
                'group' => 'system.smtp',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'smtp_username',
                'value' => null,
                'previous_value' => null,
                'group' => 'system.smtp',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'smtp_password',
                'value' => null,
                'previous_value' => null,
                'group' => 'system.smtp',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'smtp_sender_name',
                'value' => null,
                'previous_value' => null,
                'group' => 'system.smtp',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'key' => 'smtp_sender_address',
                'value' => null,
                'previous_value' => null,
                'group' => 'system.smtp',
                'created_at' => now(),
                'updated_at' => now()
            ]
        ], ['key']);
    }
}
