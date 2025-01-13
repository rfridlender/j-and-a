<script setup lang="ts">
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { toast } from "@/components/ui/toast"

import { vAutoAnimate } from "@formkit/auto-animate/vue"
import { toTypedSchema } from "@vee-validate/zod"
import { confirmSignIn, confirmSignUp, signIn, signUp } from "aws-amplify/auth"
import { LoaderCircle, LogIn, ShieldCheck } from "lucide-vue-next"
import { useForm } from "vee-validate"
import { onMounted } from "vue"
import { useRouter } from "vue-router"
import * as z from "zod"

const router = useRouter()

const signInSchema = toTypedSchema(
    z.object({
        verificationCode: z.string().min(1, "Verification code required"),
    })
)

const { handleSubmit, isSubmitting } = useForm({
    validationSchema: signInSchema,
    initialValues: {
        verificationCode: "",
    },
})

onMounted(() => {
    const successMessage = router.currentRoute.value.query["success-message"]?.toString()
    if (successMessage) {
        toast({
            title: "Success",
            description: successMessage,
        })
    }
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
                <RouterLink class="underline" to="/sign-in">Back to sign in</RouterLink>
            </div>
        </CardFooter>
    </Card>
</template>
