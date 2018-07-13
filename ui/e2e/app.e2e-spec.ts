import { DatapaceUiPage } from './app.po';

describe('datapace-ui App', () => {
  let page: DatapaceUiPage;

  beforeEach(() => {
    page = new DatapaceUiPage();
  });

  it('should display welcome message', done => {
    page.navigateTo();
    page.getParagraphText()
      .then(msg => expect(msg).toEqual('Welcome to app!!'))
      .then(done, done.fail);
  });
});
