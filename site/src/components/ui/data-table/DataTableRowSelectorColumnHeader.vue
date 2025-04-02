<script setup lang="ts" generic="TData">
import { Checkbox } from "@/components/ui/checkbox"

import type { Table } from "@tanstack/vue-table"

const props = defineProps<{ table: Table<TData> }>()
</script>

<template>
    <div class="flex items-center">
        <Checkbox
            :deselect="props.table.getIsSomeRowsSelected()"
            :checked="props.table.getIsAllPageRowsSelected()"
            @update:checked="
                (value: boolean) =>
                    value && props.table.getIsSomeRowsSelected()
                        ? props.table.toggleAllPageRowsSelected(false)
                        : props.table.toggleAllPageRowsSelected(value)
            "
            aria-label="Select all"
        ></Checkbox>
    </div>
</template>
