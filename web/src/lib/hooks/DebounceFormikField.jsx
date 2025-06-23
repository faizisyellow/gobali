import { useEffect, useRef } from "react";

export function useDebouncedFormikField(formik, fieldName, delay = 10) {
  const ref = useRef(formik.values[fieldName]);
  const timeoutRef = useRef(null);

  const onInput = (e) => {
    const val = e.target.value;
    ref.current = val;

    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = setTimeout(() => {
      if (formik.values[fieldName] !== ref.current) {
        formik.setFieldValue(fieldName, ref.current);
      }
    }, delay);
  };

  // Keep ref in sync with Formik external value change
  useEffect(() => {
    ref.current = formik.values[fieldName];
  }, [formik.values[fieldName]]);

  return {
    value: ref.current,
    onInput,
  };
}

export default useDebouncedFormikField;
