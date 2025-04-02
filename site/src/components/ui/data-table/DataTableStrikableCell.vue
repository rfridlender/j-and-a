<script setup lang="ts">
import { cn } from "@/lib/utils"
import type { ModelTypeValues } from "@/models"

import type { Row } from "@tanstack/vue-table"
import type { HTMLAttributes } from "vue"

const props = defineProps<{ class?: HTMLAttributes["class"]; row: Row<ModelTypeValues> }>()
</script>

<template>
    <div
        :class="
            cn(
                props.row.original.deletedAt && 'line-through',
                props.row.original.deletedAt && !props.row.getIsSelected() && 'opacity-50',
                props.class
            )
        "
    >
        <slot />
    </div>
</template>
