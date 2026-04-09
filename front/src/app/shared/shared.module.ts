import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { NavbarComponent } from './navbar/navbar.component';
import { CategoriesBarComponent } from './categories-bar/categories-bar.component';
import { WarningBannerComponent } from './warning-banner/warning-banner.component';
import { AuthModalComponent } from './auth-modal/auth-modal.component';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    RouterModule,
    NavbarComponent,
    CategoriesBarComponent,
    WarningBannerComponent,
    AuthModalComponent,
  ],
  exports: [
    NavbarComponent,
    CategoriesBarComponent,
    WarningBannerComponent,
    AuthModalComponent,
  ],
})
export class SharedModule { }