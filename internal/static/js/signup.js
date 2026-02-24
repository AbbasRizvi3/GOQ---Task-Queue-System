function signupForm(initialName, initialEmail) {
    return {
        formData: {
            name: initialName || '',
            email: initialEmail || '',
            password: '',
            confirmPassword: ''
        },
        errors: {},
        isSubmitting: false,
        shakeForm: false,

        validateName() {
            const regex = /^[a-zA-Z0-9 ._-]+$/;
            if (!this.formData.name.trim()) {
                this.errors.name = "Name is required";
            } else if (this.formData.name.length < 3) {
                this.errors.name = "Name must be at least 3 characters";
            } else if (this.formData.name.length > 50) {
                this.errors.name = "Name must be less than 50 characters";
            } else if (!regex.test(this.formData.name)) {
                this.errors.name = "Invalid characters used";
            } else {
                delete this.errors.name;
            }
        },

        validateEmail() {
            const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!this.formData.email.trim()) {
                this.errors.email = "Email is required";
            } else if (this.formData.email.length > 100) {
                this.errors.email = "Email is too long";
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
                this.errors.password = "Must be at least 8 characters";
            } else if (this.formData.password.length > 72) {
                this.errors.password = "Password is too long";
            } else {
                delete this.errors.password;
            }
            if (this.formData.confirmPassword) this.validateConfirm();
        },

        validateConfirm() {
            if (this.formData.password !== this.formData.confirmPassword) {
                this.errors.confirmPassword = "Passwords do not match";
            } else delete this.errors.confirmPassword;
        },

        submitForm(formElement) {
            this.validateName();
            this.validateEmail();
            this.validatePassword();
            this.validateConfirm();

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