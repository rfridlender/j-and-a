<script setup lang="ts">
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useToast } from "@/components/ui/toast"

import { vAutoAnimate } from "@formkit/auto-animate/vue"
import { toTypedSchema } from "@vee-validate/zod"
import { confirmSignIn, signIn } from "aws-amplify/auth"
import { LoaderCircle, ShieldCheck } from "lucide-vue-next"
import { useForm } from "vee-validate"
import { ref } from "vue"
import { useRouter } from "vue-router"
import * as z from "zod"

const router = useRouter()
const { toast } = useToast()

const { handleSubmit, isSubmitting } = useForm({
    validationSchema: toTypedSchema(
        z.object({
            verificationCode: z.string().min(1, "Verification code required"),
        })
    ),
    initialValues: {
        verificationCode: "",
    },
})

const onSubmit = handleSubmit(async (values) => {
    try {
        const confirmSignInOutput = await confirmSignIn({
            challengeResponse: values.verificationCode,
        })

        console.log("confirmSignInOutput", JSON.stringify(confirmSignInOutput, null, 2))

        if (confirmSignInOutput.nextStep.signInStep === "DONE") {
            router.push("/dashboard")
        }
    } catch (error) {
        console.error(error)

        if (error instanceof Error) {
            switch (error.name) {
                case "SignInException":
                    router.replace({
                        path: "/sign-in",
                        query: { "error-message": error.message },
                    })
                default:
                    toast({
                        title: error.name,
                        description: error.message,
                        variant: "destructive",
                    })
            }
        }
    }
})

const isResending = ref(false)
async function onResendVerificationCode() {
    try {
        isResending.value = true

        const email = router.currentRoute.value.query.email?.toString()
        if (!email) {
            throw new Error("Your previous session has expired")
        }

        const signInOutput = await signIn({
            username: email,
            options: {
                authFlowType: "USER_AUTH",
                preferredChallenge: "EMAIL_OTP",
            },
        })

        console.log("signInOutput", JSON.stringify(signInOutput, null, 2))

        if (signInOutput.nextStep.signInStep !== "CONFIRM_SIGN_IN_WITH_EMAIL_CODE") {
            router.push({
                path: "/sign-in",
                query: { "error-message": "Something went wrong" },
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
    } finally {
        isResending.value = false
    }
}
</script>

<template>
    <Card>
        <CardHeader>
            <CardTitle class="text-2xl">Verify</CardTitle>
            <CardDescription>Enter the verification code sent to your email</CardDescription>
        </CardHeader>

        <CardContent>
            <form id="form" @submit="onSubmit">
                <div class="grid gap-4">
                    <FormField v-slot="{ componentField }" name="verificationCode">
                        <FormItem v-auto-animate>
                            <FormLabel>Verification Code</FormLabel>

                            <FormControl>
                                <Input v-bind="componentField" type="text" placeholder="Verification Code" autocomplete="off" />
                            </FormControl>

                            <FormMessage />
                        </FormItem>
                    </FormField>
                </div>
            </form>
        </CardContent>

        <CardFooter class="flex flex-col">
            <Button class="w-full" type="submit" form="form" :disabled="isSubmitting">
                Verify
                <ShieldCheck v-if="!isSubmitting" />
                <LoaderCircle v-else class="animate-spin" />
            </Button>

            <div class="flex flex-col gap-2 mt-4 text-center text-sm">
                <div class="flex justify-center items-center gap-1">
                    <span>Didn't receive a code?</span>
                    <span class="underline cursor-pointer" @click="onResendVerificationCode">Resend</span>
                    <LoaderCircle v-if="isResending" class="size-4 animate-spin" />
                </div>

                <RouterLink class="underline" to="/sign-in">Back to sign in</RouterLink>
            </div>
        </CardFooter>
    </Card>
</template>
