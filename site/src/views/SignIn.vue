<script setup lang="ts">
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useToast } from "@/components/ui/toast"

import { vAutoAnimate } from "@formkit/auto-animate/vue"
import { toTypedSchema } from "@vee-validate/zod"
import { signIn } from "aws-amplify/auth"
import { LoaderCircle, LogIn } from "lucide-vue-next"
import { useForm } from "vee-validate"
import { onMounted } from "vue"
import { useRouter } from "vue-router"
import * as z from "zod"

const router = useRouter()
const { toast } = useToast()

const { handleSubmit, isSubmitting } = useForm({
    validationSchema: toTypedSchema(
        z.object({
            email: z.string().min(1, "Email required").email("Invalid email"),
        })
    ),
    initialValues: {
        email: "",
    },
})

onMounted(() => {
    const errorMessage = router.currentRoute.value.query["error-message"]?.toString()
    if (errorMessage) {
        toast({
            title: "Error",
            description: errorMessage,
            variant: "destructive",
        })
    }
})

const onSubmit = handleSubmit(async (values) => {
    try {
        const signInOutput = await signIn({
            username: values.email,
            options: {
                authFlowType: "USER_AUTH",
                preferredChallenge: "EMAIL_OTP",
            },
        })

        console.log("signInOutput", JSON.stringify(signInOutput, null, 2))

        if (signInOutput.nextStep.signInStep === "CONFIRM_SIGN_IN_WITH_EMAIL_CODE") {
            router.push({
                path: "/confirm-sign-in-with-email-code",
                query: { email: values.email },
            })
        }
    } catch (error) {
        console.error(error)

        if (error instanceof Error) {
            toast({
                title: error.name,
                description: error.message,
                variant: "destructive",
            })
        }
    }
})
</script>

<template>
    <Card>
        <CardHeader>
            <CardTitle class="text-2xl">Sign In</CardTitle>
            <CardDescription>Enter your email to sign in to your account</CardDescription>
        </CardHeader>

        <CardContent>
            <form id="form" @submit="onSubmit">
                <div class="grid gap-4">
                    <FormField v-slot="{ componentField }" name="email">
                        <FormItem v-auto-animate>
                            <FormLabel>Email</FormLabel>

                            <FormControl>
                                <Input v-bind="componentField" type="text" placeholder="Email" />
                            </FormControl>

                            <FormMessage />
                        </FormItem>
                    </FormField>
                </div>
            </form>
        </CardContent>

        <CardFooter class="flex flex-col">
            <Button class="w-full" type="submit" form="form" :disabled="isSubmitting">
                Sign In
                <LogIn v-if="!isSubmitting" />
                <LoaderCircle v-else class="animate-spin" />
            </Button>
        </CardFooter>
    </Card>
</template>
