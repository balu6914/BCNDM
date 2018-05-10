export class Subscription {
    constructor(
        public streamId: string,
        public userId: string,
        public pricePaid: number,
        public startDate: number,
        public endDate: number,
        public streamUrl: [string]
        ){}
}
