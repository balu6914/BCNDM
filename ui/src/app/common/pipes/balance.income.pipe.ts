import { Pipe, PipeTransform } from '@angular/core';
import { Subscription } from '../interfaces/subscription.interface';
import * as moment from 'moment';
import { TasPipe } from './converter.pipe';

@Pipe({name: 'walletBalanceStatisticPipe'})
export class WalletBalanceStatisticPipe implements PipeTransform {
  constructor(private tasPipe: TasPipe) {}
  transform(subscriptions: any[]): any[] {
    const data: any[] = [];
    subscriptions.forEach((s: any, index) => {
        data.push({
          x: moment(s.start_date).toDate(),
          y: parseFloat(this.tasPipe.transform(s.stream_price)) * s.hours,
        });
    });
    return data;
  }
}
