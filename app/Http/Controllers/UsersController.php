<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\Rules\Password;
use App\Models\User;
use Illuminate\Support\Facades\Hash;

class UsersController extends Controller
{
    //get the current user
    public function me()
    {
        $user = Auth::user();

        if (!$user) {
            return response()->json([
                'status' => 'error',
                'message' => 'Unauthorized'
            ], 401);
        }

        return response()->json([
            'status' => 'success',
            'message' => 'User fetched successfully',
            'data' => [
                'user' => $user
            ]
        ]);
    }

    //update the current user
    public function updateMe(Request $request)
    {

        $user = Auth::user();

        $validator = Validator::make($request->all(), [
            'password' => ['sometimes', 'confirmed', Password::min(8)],
            'email' => ['email', 'unique:users,email,' . $user->id],
            'name' => ['string', 'max:255'],
        ]);

        $unsetMustChangePassword = false;
        if ($request->has('password')) {
            $unsetMustChangePassword = true;
        }

        if ($validator->fails()) {
            return response()->json(
                [
                    'status' => 'error',
                    'message' => 'Validation failed',
                    'data' => [
                        'errors' => $validator->errors()
                    ]
                ],
                400
            );
        }

        try {
            $user->update($validator->validated());

            if ($unsetMustChangePassword) {
                $user->must_change_password = false;
                $user->save();
            }

            return response()->json([
                'status' => 'success',
                'message' => 'Profile updated successfully',
                'data' => [
                    'user' => $user
                ]
            ]);
        } catch (\Exception $e) {
            return response()->json(
                ['status' => 'error', 'message' => 'Failed to update profile'],
                500
            );
        }
    }


    //get all users
    public function index()
    {
        $users = User::all();

        return response()->json([
            'status' => 'success',
            'message' => 'Users fetched successfully',
            'data' => [
                'users' => $users
            ]
        ]);
    }

    //create a new user
    public function create(Request $request)
    {
        $validator = Validator::make($request->all(), [
            'password' => ['required', 'confirmed', Password::min(8)],
            'email' => ['required', 'email', 'unique:users,email'],
            'name' => ['required', 'string', 'max:255'],
            'admin' => ['boolean'],
            'active' => ['boolean'],
            'must_change_password' => ['boolean'],
        ]);

        if ($validator->fails()) {
            return response()->json(
                [
                    'status' => 'error',
                    'message' => 'Validation failed',
                    'data' => [
                        'errors' => $validator->errors()
                    ]
                ],
                400
            );
        }

        try {
            $user = User::create($validator->validated());

            return response()->json([
                'status' => 'success',
                'message' => 'User created successfully',
                'data' => [
                    'user' => $user
                ]
            ]);
        } catch (\Exception $e) {
            return response()->json(
                ['status' => 'error', 'message' => 'Failed to create user'],
                500
            );
        }
    }

    //update a user
    public function update(Request $request, $id)
    {
        $user = User::find($id);

        if (!$user) {
            return response()->json(
                ['status' => 'error', 'message' => 'User not found'],
                404
            );
        }

        $validator = Validator::make($request->all(), [
            'password' => ['confirmed', Password::min(8)],
            'email' => ['email', 'unique:users,email,' . $user->id],
            'name' => ['string', 'max:255'],
            'must_change_password' => ['boolean'],
            'admin' => ['boolean'],
        ]);

        if ($validator->fails()) {
            return response()->json(
                [
                    'status' => 'error',
                    'message' => 'Validation failed',
                    'data' => [
                        'errors' => $validator->errors()
                    ]
                ],
                400
            );
        }

        try {
            $user->update($validator->validated());

            return response()->json([
                'status' => 'success',
                'message' => 'User updated successfully',
                'data' => [
                    'user' => $user
                ]
            ]);
        } catch (\Exception $e) {
            return response()->json(
                ['status' => 'error', 'message' => 'Failed to update user'],
                500
            );
        }
    }

    //delete a user
    public function delete($id)
    {
        $user = User::find($id);

        if (!$user) {
            return response()->json(
                ['status' => 'error', 'message' => 'User not found'],
                404
            );
        }

        try {
            $user->delete();

            return response()->json([
                'status' => 'success',
                'message' => 'User deleted successfully'
            ]);
        } catch (\Exception $e) {
            return response()->json(
                ['status' => 'error', 'message' => 'Failed to delete user'],
                500
            );
        }
    }


    //create the first user
    public function createFirstUser(Request $request)
    {
        $validator = Validator::make($request->all(), [
            'password' => ['required', 'confirmed', Password::min(8)],
            'email' => ['required', 'email', 'unique:users,email'],
            'name' => ['required', 'string', 'max:255'],
        ]);

        if ($validator->fails()) {
            return response()->json(
                [
                    'status' => 'error',
                    'message' => 'Validation failed',
                    'data' => [
                        'errors' => $validator->errors()

                    ]
                ],
                400
            );
        }

        try {
            $user = User::create([
                'name' => $request->name,
                'email' => $request->email,
                'password' => Hash::make($request->password),
                'admin' => true,
                'active' => true,
                'must_change_password' => false,
            ]);

            return response()->json([
                'status' => 'success',
                'message' => 'First user created successfully',
                'data' => [
                    'user' => $user
                ]
            ]);
        } catch (\Exception $e) {
            return response()->json(
                ['status' => 'error', 'message' => 'Failed to create first user'],
                500
            );
        }
    }
}
