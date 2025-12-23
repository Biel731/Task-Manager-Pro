import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  templateUrl: './register.html',
  styleUrls: ['./register.scss'],
})
export class RegisterComponent {
  loading = false;
  errorMsg = '';

  form: ReturnType<FormBuilder['group']>;

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {
    this.form = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  get name() {
    return this.form.get('name');
  }

  get email() {
    return this.form.get('email');
  }

  get password() {
    return this.form.get('password');
  }

  submit(): void {
    this.errorMsg = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const payload = {
      name: String(this.name?.value || '').trim(),
      email: String(this.email?.value || '').trim(),
      password: String(this.password?.value || ''),
    };

    this.loading = true;

    this.auth.register(payload).subscribe({
      next: () => {
        this.loading = false;
        // padrão seguro: cria conta e manda logar
        this.router.navigateByUrl('/login');
      },
      error: (err) => {
        this.loading = false;
        console.error('REGISTER ERROR:', err);
        this.errorMsg =
          err?.error?.error ||
          err?.error?.message ||
          err?.message ||
          'Falha ao criar conta';
      },
    });
  }
}
