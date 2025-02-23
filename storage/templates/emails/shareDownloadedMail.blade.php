@extends('emails.layout')

@section('header', 'Share Downloaded')
@section('content')
  <h2>Hi {{ $share->user->name }},</h2>
  <p>Just thought you'd like to know that your share "{{ $share->name }}" was downloaded.</p>
  <p><small>We'll only send this the first time a share is downloaded.</small></p>

@endsection
@section('action_url', $settings['application_url'] . '/shares/' . $share->long_id)
@section('action_text', 'View Share')
