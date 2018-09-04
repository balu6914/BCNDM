import { Pipe, PipeTransform } from '@angular/core';
import { Subscription } from '../interfaces/subscription.interface';

@Pipe({name: 'walletBalanceStatisticPipe', pure: false})
export class WalletBalanceStatisticPipe implements PipeTransform {
  transform(subscriptions: any[]): any[] {
    const data: any[] = [];
    subscriptions.forEach((s: any) => {
        data.push({
          x: new Date(),
          y: 10
        });
    });
    return data;
  }
}
