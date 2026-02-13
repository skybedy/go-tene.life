<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class CameraUploadController extends Controller
{
    /**
     * Upload image from AXIS camera
     */
    public function upload(Request $request)
    {
        // Debug - logujeme do souboru
        file_put_contents('/tmp/camera_debug.log', json_encode([
            'time' => date('Y-m-d H:i:s'),
            'method' => $request->method(),
            'content_type' => $request->header('Content-Type'),
            'has_files' => count($request->allFiles()) > 0,
            'files' => array_keys($request->allFiles()),
            'all_data_keys' => array_keys($request->all()),
            'body_length' => strlen($request->getContent()),
            'body_preview' => substr($request->getContent(), 0, 200)
        ], JSON_PRETTY_PRINT) . "\n\n", FILE_APPEND);

        // Autentizace pomocí HTTP Basic Auth - DOČASNĚ VYPNUTO PRO DEBUG
        // $username = $request->getUser();
        // $password = $request->getPassword();

        // $expectedUsername = env('CAMERA_USERNAME');
        // $expectedPassword = env('CAMERA_PASSWORD');

        // if (!$username || !$password ||
        //     $username !== $expectedUsername ||
        //     $password !== $expectedPassword) {
        //     return response()->json([
        //         'success' => false,
        //         'error' => 'Unauthorized'
        //     ], 401);
        // }

        try {
            // Název souboru - vždy stejný, přepíše předchozí
            $filename = 'tenelife.jpg';

            // Cesta k uložení
            $directory = 'images';
            $publicPath = public_path($directory);

            // Vytvoření složky pokud neexistuje
            if (!file_exists($publicPath)) {
                mkdir($publicPath, 0755, true);
            }

            // Zkusíme multipart file
            if ($request->hasFile('image')) {
                $request->file('image')->move($publicPath, $filename);
            }
            // Nebo raw body
            elseif (strlen($request->getContent()) > 0) {
                file_put_contents($publicPath . '/' . $filename, $request->getContent());
            }
            // Nebo zkusíme všechny files
            elseif (count($request->allFiles()) > 0) {
                $file = array_values($request->allFiles())[0];
                $file->move($publicPath, $filename);
            }
            else {
                throw new \Exception('No image data received');
            }

            return response()->json([
                'success' => true,
                'filename' => $filename,
                'path' => $directory . '/' . $filename,
                'url' => url($directory . '/' . $filename)
            ], 200);

        } catch (\Exception $e) {
            return response()->json([
                'success' => false,
                'error' => 'Failed to save image',
                'message' => $e->getMessage()
            ], 500);
        }
    }
}
