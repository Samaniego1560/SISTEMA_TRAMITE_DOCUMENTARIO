import {AbstractControl, FormGroup, ValidationErrors, ValidatorFn, Validators} from '@angular/forms';
import {VALIDATOR_MESSAGE_DEFAULT} from "../constants/validator-messages.constant";
import {PatternAlphabeticEmail} from "../constants/patterns.constant";

export class NgxValidators {
  /**
   * A static method that returns a validator function for required validation.
   *
   * @param {string} [message] - Optional message to be displayed when the validation fails.
   *
   * @return {ValidatorFn} - The validator function to be used for required validation.
   */
  static required(message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const error = Validators.required(control);
      return error ? {required: this._getMessage('required', message)} : null;
    };
  }

  /**
   * Validates if the input value is a valid email address.
   *
   * @param {string} [message] - Optional custom error message.
   *
   * @return {ValidatorFn} A validator function that returns a validation error object if the email is invalid,
   * or `null` if the email is valid.
   *
   * @example
   * email() // Returns a validator function that checks if the input is a valid email address.
   *
   * @example
   * email('Invalid email') // Returns a validator function that checks if the input is a valid email address,
   * with a custom error message set to 'Invalid email' if the email is invalid.
   *
   */
  static email(message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const emailError = Validators.email(control);
      const patternError = !PatternAlphabeticEmail.test(control.value);

      if (emailError) {
        return { email: this._getMessage('email', message) };
      }

      if (patternError) {
        return { email: this._getMessage('email', message) };
      }

      return null;
    };
  }

  /**
   * Returns a validator function that checks if the value of the control is greater than or equal to the specified minimum value.
   *
   * @param {number} min - The minimum value that the control's value should be greater than or equal to.
   * @param {string} [message] - Optional custom error message to be displayed when the validation fails.
   *
   * @returns {ValidatorFn} The validator function.
   *
   * @example
   * const myValidator = MyValidators.min(5, 'Value must be at least 5');
   * const control = new FormControl(3);
   * const errors = myValidator(control); // { min: 'Value must be at least 5' }
   */
  static min(min: number, message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const minFunction = Validators.min(min);
      const error = minFunction(control);

      return error ? {
        min: this._getMessage('min', message, [{
          min: min,
          current: control?.value || 0
        }])
      } : null;
    };
  }

  /**
   * Creates a validator function that validates whether the value of the control is less than or equal to a maximum value.
   *
   * @param {number} max - The maximum value allowed.
   * @param {string} [message] - An optional error message to be displayed.
   * @returns {ValidatorFn} A validator function that returns a ValidationErrors object if the control value exceeds the maximum value, or null if it is valid.
   */
  static max(max: number, message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const maxFunction = Validators.max(max);
      const error = maxFunction(control);

      return error ? {max: this._getMessage('max', message, [{max: max, current: control?.value || 0}])} : null;
    };
  }

  /**
   * Validates the minimum length of a control's value.
   *
   * @param minLength - The minimum length value.
   * @param message - Optional. The custom error message to display.
   *
   * @returns A `ValidatorFn` that validates the minimum length of the control's value.
   */
  static minLength(minLength: number, message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const minLengthFunction = Validators.minLength(minLength);
      const error = minLengthFunction(control);

      return error ? {
        minLength: this._getMessage('minLength', message, [{
          minLength: minLength,
          current: (control?.value || '0').length
        }])
      } : null;
    };
  }

  /**
   * Creates a validator function that checks if the value of a form control does not exceed a maximum length.
   *
   * @param {number} maxLength - The maximum length the value can be.
   * @param {string} [message] - Optional custom error message.
   * @returns {ValidatorFn} - The validator function.
   */
  static maxLength(maxLength: number, message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const maxLengthFunction = Validators.maxLength(maxLength);
      const error = maxLengthFunction(control);

      return error ? {
        maxLength: this._getMessage('maxLength', message, [{
          maxLength: maxLength,
          current: (control?.value || '0').length
        }])
      } : null;
    };
  }

  public static dateRange(startDateField: string, endDateField: string, message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const group = control.parent as FormGroup;
      if (!group) return null;

      const startDate = group.get(startDateField)?.value;
      const endDate = group.get(endDateField)?.value;

      if (!startDate || !endDate) return null

      return new Date(startDate).getTime() > new Date(endDate).getTime()
        ? {
          dateRange: this._getMessage('dateRange', message, [{
            dateRange: endDate,
            current: control?.value || ''
          }])
        }
        : null;
    };
  }

  private static _getMessage(
    control: keyof typeof VALIDATOR_MESSAGE_DEFAULT,
    message?: string,
    paramsMessage?: { [key: string]: unknown }[]
  ) {
    if (message) return message;

    let messageControl = VALIDATOR_MESSAGE_DEFAULT[control];
    const existParams = paramsMessage && paramsMessage.length > 0;

    if (existParams) {
      paramsMessage.forEach((params) => {
        Object.keys(params)
          .filter((key) => params[key])
          .forEach((key) => {
            messageControl = messageControl.replace(`\${${key}}`, params[key]!.toString());
          });
      });

      return messageControl;
    }

    return messageControl;
  }

  static pattern(pattern: string | RegExp, message?: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const patternFunction = Validators.pattern(pattern);
      const error = patternFunction(control);
      return error ? { pattern: this._getMessage('pattern', message) } : null;
    };
  }


}
