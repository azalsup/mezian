import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-ads-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './ads-list.component.html',
  styleUrls: ['./ads-list.component.scss'],
})
export class AdsListComponent {
  @Input() category = 'Selected';

  ads = [
    { title: 'Renovated apartment in Casablanca', detail: '3 bedrooms • 120m² • 1 750 000 MAD', badge: 'Real estate' },
    { title: '2018 Renault Clio, very good condition', detail: 'Diesel • 85 000 km • 125 000 MAD', badge: 'Vehicles' },
    { title: 'Freelance graphic designer', detail: 'Remote • Starting at 2 000 MAD', badge: 'Jobs' },
  ];
}
