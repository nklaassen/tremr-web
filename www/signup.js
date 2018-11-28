//	Name of file: signup.js
//	Programmers: Colin Chan and Devansh Chopra and Nic Klaassen
//	Team Name: Co.DEsign
//	Changes been made:
//					2018-11-26: created file
// Known Bugs:

password = document.getElementById('password')
confirm_password = document.getElementById('confirm_password')
full_name = document.getElementById('full_name')
email = document.getElementById('email')

function validatePassword() {
	if (password.value != confirm_password.value) {
		confirm_password.setCustomValidity('Passwords do not match')
	} else {
		confirm_password.setCustomValidity('')
	}
}

password.onchange = validatePassword
confirm_password.onkeyup = validatePassword

function onSignup() {
	formData = {
		'name': full_name.value,
		'email': email.value,
		'password': password.value
	}

	fetch('/api/auth/signup', {
		method: 'POST',
		body: JSON.stringify(formData)
	}).then(response => {
		if (response.status != 200) {
			response.text().then(errorText => alert(errorText))
		} else {
			window.location.replace('/signin.html')
		}
		return false
	}).catch(error => {
		alert(error)
	})
	return false
}
