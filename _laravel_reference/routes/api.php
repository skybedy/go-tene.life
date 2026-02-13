<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\CameraUploadController;

// API endpoint for camera upload (no CSRF protection, no redirects)
// Support POST, PUT and GET for compatibility with AXIS cameras
Route::match(['post', 'put', 'get'], '/camera/upload', [CameraUploadController::class, 'upload'])->name('api.camera.upload');
