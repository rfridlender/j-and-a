<script setup lang="ts" generic="TData">
import { Checkbox } from "@/components/ui/checkbox"

import type { Row } from "@tanstack/vue-table"

const props = defineProps<{ row: Row<TData> }>()
</script>

<template>
    <div class="flex items-center">
        <Checkbox
            :checked="props.row.getIsSelected()"
            @update:checked="(value: boolean) => props.row.toggleSelected(!!value)"
            aria-label="Select row"
        ></Checkbox>
    </div>
</template>
