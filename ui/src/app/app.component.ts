import { Component } from '@angular/core';
import { MdlDialogService, MdlDialogOutletService } from '@angular-mdl/core';
import { ViewContainerRef } from '@angular/core';

@Component({
  selector: 'dpc-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'Datapace';

  constructor(private vcRef: ViewContainerRef, private dialogService: MdlDialogOutletService) {
        this.dialogService.setDefaultViewContainerRef(this.vcRef);
  }
}
