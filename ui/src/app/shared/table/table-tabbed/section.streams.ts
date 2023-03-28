export enum StreamsType {
  Bought = 'Bought',
  Sold = 'Sold',
  All = 'All'
}

export class StreamSection {
  name: StreamsType = StreamsType.Bought;
}
