import { Component, inject, signal, OnInit } from '@angular/core';
import { LangService } from '../../core/services/lang.service';
import { BackendService } from '../../core/services/backend.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-warning-banner',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './warning-banner.component.html',
  styleUrl: './warning-banner.component.scss'
})
export class WarningBannerComponent implements OnInit {
  protected readonly lang = inject(LangService);
  private readonly backend = inject(BackendService);

  readonly message = signal<string>(this.lang.t('testingWarning'));

  ngOnInit(): void {
    console.log('WarningBanner ngOnInit');
    this.checkHealth();
  }

  private checkHealth(): void {
    // TODO: Disable test message later
    let msg = this.lang.t('testingWarning');

    console.log('Checking health...');
    this.backend.checkHealth().subscribe({
      next: (isHealthy) => {
        console.log('Health check result:', isHealthy);
        if (!isHealthy) {
          // Backend down, cumulate with maintenance
          msg += ' | ' + this.lang.t('maintenanceWarning');
        }
        this.message.set(msg);
      }
    });
  }
}
