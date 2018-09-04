import { Pipe, PipeTransform } from '@angular/core';
import { Subscription } from '../interfaces/subscription.interface';
import * as moment from 'moment';

@Pipe({name: 'walletBalanceStatisticPipe', pure: false})
export class WalletBalanceStatisticPipe implements PipeTransform {
  transform(subscriptions: any[]): any[] {
    const data: any[] = [];
    subscriptions.forEach((s: any) => {
        data.push({
          x:  moment(s.start_date).toDate(),
          y: 20
        });
    });
    return data;
  }
}
