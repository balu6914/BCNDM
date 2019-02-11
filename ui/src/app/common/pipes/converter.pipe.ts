import { Pipe, PipeTransform } from '@angular/core';

// DPC to miDPC coefficient (1 DPC = 10^6 miDPC)
const dpcToMidpcCoef = Math.pow(10, 6);

@Pipe({name: 'mipdc'})
export class MidpcPipe implements PipeTransform {
  transform(value: any): any {
      // Convert DPC to miDPC
      const mipdc = parseFloat(value) * dpcToMidpcCoef;
      return mipdc;
  }
}

@Pipe({name: 'dpc'})
export class DpcPipe implements PipeTransform {
  transform(value: any): any {
      // Convert miDPC to DPC
      const dpc = parseInt(value) / dpcToMidpcCoef;
      return dpc.toString();
  }
}

@Pipe({name: 'subscriptionType'})
export class SubscriptionTypePipe implements PipeTransform {
  transform(value: number, type: string): any {
    if (type === 'Income') {
      return '+' + value;
    }
    if (type === 'Outcome') {
      return '-' + value;
    }
    return value.toString();
  }
}
