<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Download extends Model
{
    protected $fillable = [
        'share_id',
        'user_id',
        'ip_address',
        'user_agent',
    ];
    
    public function share()
    {
        return $this->belongsTo(Share::class);
    }

    public function user()
    {
        return $this->belongsTo(User::class);
    }

    public function scopeByShare($query, $share)
    {
        return $query->where('share_id', $share->id);
    }

    public function scopeByUser($query, $user)
    {
        return $query->where('user_id', $user->id);
    }

    public function scopeByIpAddress($query, $ipAddress)
    {
        return $query->where('ip_address', $ipAddress);
    }

    public function scopeByUserAgent($query, $userAgent)
    {
        return $query->where('user_agent', $userAgent);
    }
    
}
