//	Name of file: signin.js
//	Programmers: Colin Chan and Devansh Chopra and Nic Klaassen
//	Team Name: Co.DEsign
//	Changes been made:
//					2018-11-26: created file
// Known Bugs:

function onSignin() {
	email = document.getElementById('email').value
	password = document.getElementById('password').value

	formData = {
		'email': email,
		'password': password
	}

	fetch('/api/auth/signin', {
		method: 'POST',
		body: JSON.stringify(formData)
	}).then(response => {
		if (response.status != 200) {
			response.text().then(errorText => alert(errorText))
		} else {
			response.text().then(token => {
				localStorage.setItem('token', token)
				window.location.replace('/index.html')
			})
		}
		return false
	}).catch(error => {
		alert(error)
	})
	return false
}
