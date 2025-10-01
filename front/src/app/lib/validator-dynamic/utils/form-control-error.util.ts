import {AbstractControl} from "@angular/forms";
import {VALIDATOR_MESSAGE_DEFAULT} from "../constants/validator-messages.constant";

export const getFormControlError = (formControl: AbstractControl): string => {
  if (!formControl.errors) return '';

  const firstErrorKey = Object.keys(formControl.errors!)[0];

  if (formControl.errors[firstErrorKey] === true) {
    return VALIDATOR_MESSAGE_DEFAULT[firstErrorKey];
  }

  return formControl.errors![firstErrorKey] || '';
};
