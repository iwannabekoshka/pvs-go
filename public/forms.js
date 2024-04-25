const forms = document.querySelectorAll(".js-validation");

forms.forEach(form => form.addEventListener("submit", formSubmitHandler));

async function formSubmitHandler(event) {
  event.preventDefault();

  const form = event.target;

  const formData = new FormData(form);

  const response = await fetch(form.action, {
    method: 'POST',
    body: formData
  });

  const body = await response.json();

  form.querySelectorAll(".error").forEach(error => error.remove())

  if (body.errors && body.errors.length > 0) {
    for (const error of body.errors) {
      const p = document.createElement("p");
      p.classList.add("error");
      p.textContent = error.Message;

      const input = form.querySelector(`[name="${error.Column.toLowerCase()}"]`);
      input.after(p)
    }
  } else {
    const p = document.createElement("p");
    p.classList.add("success");
    p.textContent = body.message;
    form.after(p);
  }
}