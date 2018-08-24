import { Pipe, PipeTransform } from '@angular/core';

// TAS to miTAS coefficient (1 TAS = 10^6 miTAS)
const tasToMitasCoef = Math.pow(10, 6);

@Pipe({name: 'mitas'})
export class MitasPipe implements PipeTransform {
  transform(value: any): any {
      // Convert TAS to miTAS
      const mitas = parseFloat(value) * tasToMitasCoef;
      return mitas;
  }
}

@Pipe({name: 'tas'})
export class TasPipe implements PipeTransform {
  transform(value: any): any {
      // Convert miTAS to TAS
      const tas = parseInt(value) / tasToMitasCoef;
      return tas.toString();
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
