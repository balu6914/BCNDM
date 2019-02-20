import { Pipe, PipeTransform } from '@angular/core';
import { Subscription } from 'app/common/interfaces/subscription.interface';
import * as moment from 'moment';
import { DpcPipe } from './converter.pipe';

@Pipe({name: 'walletBalanceStatisticPipe'})
export class WalletBalanceStatisticPipe implements PipeTransform {
  constructor(private dpcPipe: DpcPipe) {}
  transform(subscriptions: any[]): any[] {
    const data: any[] = [];
    subscriptions.forEach((s: any, index) => {
        data.push({
          x: moment(s.start_date).toDate(),
          y: parseFloat(this.dpcPipe.transform(s.stream_price)) * s.hours,
        });
    });
    return data;
  }
}
