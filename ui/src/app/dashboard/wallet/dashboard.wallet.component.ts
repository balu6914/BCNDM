import { Component, OnInit } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';

import { AuthService } from 'app/auth/services/auth.service';
import { User } from 'app/common/interfaces/user.interface';

@Component({
  selector: 'user-wallet-balance',
  templateUrl: './dashboard.wallet.component.html',
  styleUrls: ['./dashboard.wallet.component.scss']
})
export class DashboardWalletComponent implements OnInit {
  user: User;
  newBalance: number;
  subscription: any;
  // TODO: Remove this Mock of user balance its tmp hack for balance wallet widget
  mockBalance: any;

  constructor(
    private authService: AuthService,
    private modalService: BsModalService,
  ) {}

    ngOnInit() {
      this.authService.getCurrentUser().subscribe(
        data =>  {
          this.user = data
        },
        err => {
          console.log(err);
        }
      );
    }
  }
