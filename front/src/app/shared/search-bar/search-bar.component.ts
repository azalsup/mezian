import { Component, Input, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';

@Component({
  selector: 'app-search-bar',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './search-bar.component.html',
  styleUrls: ['./search-bar.component.scss'],
})
export class SearchBarComponent {
  protected readonly lang = inject(LangService);
  @Input() categories: Array<{ icon: string; slug: string; labelFr: string; labelAr: string; count: string; subcategories: string[] }> = [];
}
