import { BlocksenseUiPage } from './app.po';

describe('blocksense-ui App', () => {
  let page: BlocksenseUiPage;

  beforeEach(() => {
    page = new BlocksenseUiPage();
  });

  it('should display welcome message', done => {
    page.navigateTo();
    page.getParagraphText()
      .then(msg => expect(msg).toEqual('Welcome to app!!'))
      .then(done, done.fail);
  });
});
