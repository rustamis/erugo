<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::table('shares', function (Blueprint $table) {
            $table->boolean('sent_expiry_warning')->default(false);
            $table->boolean('sent_expired')->default(false);
            $table->boolean('sent_deletion_warning')->default(false);
            $table->boolean('sent_deleted')->default(false);
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('shares', function (Blueprint $table) {
            $table->dropColumn('sent_expiry_warning');
            $table->dropColumn('sent_expired');
            $table->dropColumn('sent_deletion_warning');
            $table->dropColumn('sent_deleted');
        });
    }
};
