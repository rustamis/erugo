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
        Schema::create('shares', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->constrained('users');
            $table->string('name')->nullable();
            $table->string('description')->nullable();
            $table->string('path');
            $table->string('password')->nullable();
            $table->string('long_id')->unique();
            $table->integer('size');
            $table->integer('file_count');
            $table->integer('download_limit')->nullable();
            $table->integer('download_count')->default(0);
            $table->boolean('require_email')->default(false);
            
            $table->timestamps();
            $table->dateTime('expires_at')->nullable();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('shares');
    }
};
