<script setup lang="ts">
import banner from "@/assets/banner.png"

import { RouterView } from "vue-router"
</script>

<template>
    <main class="w-full h-full lg:grid lg:grid-cols-2">
        <div class="w-full flex justify-center items-center">
            <div class="max-w-sm">
                <div class="p-6">
                    <img :src="banner" alt="Placeholder" />
                </div>

                <RouterView />
            </div>
        </div>

        <div class="hidden lg:block relative bg-muted">
            <div class="w-full h-full absolute top-0 left-0 bg-black/50" />
            <img class="w-full h-full object-cover" :src="banner" alt="Placeholder" />
        </div>
    </main>
</template>
