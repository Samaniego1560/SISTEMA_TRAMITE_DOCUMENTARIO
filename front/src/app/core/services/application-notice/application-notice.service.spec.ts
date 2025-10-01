import { TestBed } from '@angular/core/testing';

import { ApplicationNoticeService } from './application-notice.service';

describe('ApplicationNoticeService', () => {
  let service: ApplicationNoticeService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ApplicationNoticeService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
