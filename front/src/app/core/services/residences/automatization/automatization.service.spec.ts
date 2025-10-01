import { TestBed } from '@angular/core/testing';

import { AutomatizationService } from './automatization.service';

describe('AutomatizationService', () => {
  let service: AutomatizationService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AutomatizationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
