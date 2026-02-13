<x-app-layout>
    <div class="mx-auto bg-white/30 overflow-hidden w-full border border-gray-300">

        <div class="relative">

            @php
                $currentLocale = app()->getLocale();
            @endphp

            <a href="{{ url(($currentLocale && $currentLocale !== 'cs' ? '/' . $currentLocale : '') . '/') }}"
               class="absolute top-3 lg:top-5 right-3 lg:right-5
                      bg-black/40
                      hover:bg-white/50
                      text-white
                      hover:text-black
                      pl-1 pr-3 sm:px-4 py-1 sm:py-2
                      rounded-full
                      backdrop-blur-sm shadow-lg
                      flex items-center gap-2 transition">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M15 19l-7-7 7-7" />
                </svg>
                <span>{{ __('messages.back') }}</span>
            </a>
            
            <img src="{{ asset('images/tenelife.jpg') }}" alt="Webkamera výhled" class="w-full">

           
        </div>

    </div>
</x-app-layout>




