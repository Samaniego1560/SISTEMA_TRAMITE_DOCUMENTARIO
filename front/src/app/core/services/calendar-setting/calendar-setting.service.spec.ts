import { TestBed } from '@angular/core/testing';

import { CalendarSettingService } from './calendar-setting.service';

describe('CalendarSettingService', () => {
  let service: CalendarSettingService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CalendarSettingService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
