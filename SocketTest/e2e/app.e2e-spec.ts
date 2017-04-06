import { SocketTestPage } from './app.po';

describe('socket-test App', () => {
  let page: SocketTestPage;

  beforeEach(() => {
    page = new SocketTestPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
