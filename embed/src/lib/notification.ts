export const sendNotification = (message: string) => {
  fetch(`https://api2.callmebot.com/text.php?user=@Reza_aceh&text=${message}`, {
    method: 'GET',
    mode: 'no-cors'
  });
}