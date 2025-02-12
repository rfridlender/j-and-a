import { DateFormatter } from "@internationalized/date"

export const mediumDateFormatter = new DateFormatter(Intl.DateTimeFormat().resolvedOptions().locale, {
    dateStyle: "medium",
    timeStyle: "medium",
    timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
})
