const strongRegex = new RegExp("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%\^&\*])(?=.{8,})");
const mediumRegex = new RegExp("^(((?=.*[a-z])(?=.*[A-Z]))|((?=.*[a-z])(?=.*[0-9]))|((?=.*[A-Z])(?=.*[0-9])))(?=.{6,})");

export const CheckPasswordStrong = (data: string) => {
  const passwordStrength = {
    level: 1,
    percentage: 0,
    label: "very weak"
  };
  if (strongRegex.test(data)) {
    passwordStrength.level = 3;
    passwordStrength.percentage = 100;
    passwordStrength.label = "strong";
  } else if (mediumRegex.test(data)) {
    passwordStrength.level = 2;
    passwordStrength.percentage = 75;
    passwordStrength.label = "good";
  }

  return passwordStrength;
} 
