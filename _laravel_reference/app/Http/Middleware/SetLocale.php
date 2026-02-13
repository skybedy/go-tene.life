<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Session;
use Symfony\Component\HttpFoundation\Response;

class SetLocale
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        $supportedLocales = ['cs', 'en', 'es', 'de', 'it', 'pl', 'hu', 'fr'];
        $locale = $request->segment(1);

        // If locale is in URL and is supported, use it
        if (in_array($locale, $supportedLocales)) {
            App::setLocale($locale);
            Session::put('locale', $locale);
        } else {
            // No locale prefix means Czech (default)
            App::setLocale('cs');
            Session::put('locale', 'cs');
        }

        return $next($request);
    }
}
