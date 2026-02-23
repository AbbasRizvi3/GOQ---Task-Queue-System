function loginForm() {
    return {
        formData: {
            email: '',
            password: ''
        },
        errors: {},
        isSubmitting: false,
        shakeForm: false,

        validateEmail() {
            const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!this.formData.email.trim()) {
                this.errors.email = "Email is required";
            } else if (!re.test(this.formData.email)) {
                this.errors.email = "Invalid email address";
            } else {
                delete this.errors.email;
            }
        },

        validatePassword() {
            if (!this.formData.password) {
                this.errors.password = "Password is required";
            } else if (this.formData.password.length < 8) {
                this.errors.password = "Password must be at least 8 characters";
            } else {
                delete this.errors.password;
            }
        },

        submitForm(formElement) {
            this.validateEmail();
            this.validatePassword();

            if (Object.keys(this.errors).length > 0) {
                this.shakeForm = true;
                setTimeout(() => this.shakeForm = false, 500);
                return;
            }

            this.isSubmitting = true;
            formElement.submit();
        }
    }
}