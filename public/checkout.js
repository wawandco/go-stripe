// This is your test publishable API key.
const stripe = Stripe("pk_test_51IIiV0C5e5WNMZdtXlqQwEJi8CxGjKxJyz5n3Cz0dpmAjrRoYEWMkJdJdKyFUBM61jErz3nR356MHMMKiQjnCQqL00XdOj4C8G");

// The items the customer wants to buy
const items = [{ id: "xl-tshirt" }];

let elements;

checkStatus();

document
  .querySelector("#payment-form")
  .addEventListener("submit", handleSubmit);

// Fetches a payment intent and captures the client secret
async function startPaymentIntent() {
  const response = await fetch("/create-payment-intent/", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ items }),
  });
  const { clientSecret } = await response.json();
  const appearance = {
    theme: 'stripe',
  };
  elements = stripe.elements({ appearance, clientSecret });

  const paymentElementOptions = {
    layout: "tabs",
  };

  const paymentForm = document.querySelector("#payment-form");
  paymentForm.classList.remove("hidden");
  const paymentElement = elements.create("payment", paymentElementOptions);
  paymentElement.mount("#payment-element");

  const nextButton = document.querySelector("#next-button");
  nextButton.classList.add("hidden");
  
}

async function handleSubmit(e) {
  e.preventDefault();
  setLoading(true);

  const { error } = await stripe.confirmPayment({
    elements,
    confirmParams: {
      // Make sure to change this to your payment completion page
      return_url: "http://localhost:3000#page-3",
    },
  });

  // This point will only be reached if there is an immediate error when
  // confirming the payment. Otherwise, your customer will be redirected to
  // your `return_url`. For some payment methods like iDEAL, your customer will
  // be redirected to an intermediate site first to authorize the payment, then
  // redirected to the `return_url`.
  if (error.type === "card_error" || error.type === "validation_error") {
    showMessage(error.message);
  } else {
    showMessage("An unexpected error occurred.");
  }

  setLoading(false);
}

// Fetches the payment intent status after payment submission
async function checkStatus() {
  const clientSecret = new URLSearchParams(window.location.search).get(
    "payment_intent_client_secret"
  );

  if (!clientSecret) {
    return;
  }

  const { paymentIntent } = await stripe.retrievePaymentIntent(clientSecret);

  switch (paymentIntent.status) {
    case "succeeded":
      showMessage("Payment succeeded!");
      
      break;
    case "processing":
      showMessage("Your payment is processing.");
      break;
    case "requires_payment_method":
      showMessage("Your payment was not successful, please try again.");
      break;
    default:
      showMessage("Something went wrong.");
      break;
  }
}

// ------- UI helpers -------

function showMessage(messageText, selector="#payment-message") {
  const messageContainer = document.querySelector(selector);

  messageContainer.classList.remove("hidden");
  messageContainer.textContent = messageText;

  setTimeout(function () {
    messageContainer.classList.add("hidden");
    messageContainer.textContent = "";
  }, 4000);
}

// Show a spinner on payment submission
function setLoading(isLoading) {
  if (isLoading) {
    // Disable the button and show a spinner
    document.querySelector("#submit").disabled = true;
    document.querySelector("#spinner").classList.remove("hidden");
    document.querySelector("#button-text").classList.add("hidden");
  } else {
    document.querySelector("#submit").disabled = false;
    document.querySelector("#spinner").classList.add("hidden");
    document.querySelector("#button-text").classList.remove("hidden");
  }
}


// -------

checkStatusAppFee();
document
.querySelector("#payment-form-app-fee")
.addEventListener("submit", handleSubmitAppFee);


// Fetches a payment intent and captures the client secret
async function startPaymentIntentAppFee() {
  const response = await fetch("/create-payment-intent-app-fee/", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ items }),
  });
  const { clientSecret } = await response.json();
  const appearance = {
    theme: 'stripe',
  };
  elements = stripe.elements({ appearance, clientSecret });

  const paymentElementOptions = {
    layout: "tabs",
  };

  const paymentForm = document.querySelector("#payment-form-app-fee");
  paymentForm.classList.remove("hidden");
  const paymentElement = elements.create("payment", paymentElementOptions);
  paymentElement.mount("#payment-element-app-fee");

  const nextButton = document.querySelector("#next-button-app-fee");
  nextButton.classList.add("hidden");
  
}

async function handleSubmitAppFee(e) {
  e.preventDefault();
  setLoading(true);

  const { error } = await stripe.confirmPayment({
    elements,
    confirmParams: {
      // Make sure to change this to your payment completion page
      return_url: "http://localhost:3000#page-5",
    },
  });

  // This point will only be reached if there is an immediate error when
  // confirming the payment. Otherwise, your customer will be redirected to
  // your `return_url`. For some payment methods like iDEAL, your customer will
  // be redirected to an intermediate site first to authorize the payment, then
  // redirected to the `return_url`.
  if (error.type === "card_error" || error.type === "validation_error") {
    showMessage(error.message, "#payment-message-app-fee");
  } else {
    showMessage("An unexpected error occurred.", "#payment-message-app-fee");
  }

  setLoading(false);
}

async function checkStatusAppFee() {
  const clientSecret = new URLSearchParams(window.location.search).get(
    "payment_intent_client_secret"
  );

  if (!clientSecret) {
    return;
  }

  const { paymentIntent } = await stripe.retrievePaymentIntent(clientSecret);

  switch (paymentIntent.status) {
    case "succeeded":
      showMessage("Payment succeeded!", "#payment-message-app-fee");
      
      break;
    case "processing":
      showMessage("Your payment is processing.", "#payment-message-app-fee");
      break;
    case "requires_payment_method":
      showMessage("Your payment was not successful, please try again.", "#payment-message-app-fee");
      break;
    default:
      showMessage("Something went wrong.", "#payment-message-app-fee");
      break;
  }
}

function setLoadingAppFee(isLoading) {
  if (isLoading) {
    // Disable the button and show a spinner
    document.querySelector("#submit-app-fee").disabled = true;
    document.querySelector("#spinner-app-fee").classList.remove("hidden");
    document.querySelector("#button-text-app-fee").classList.add("hidden");
  } else {
    document.querySelector("#submit-app-fee").disabled = false;
    document.querySelector("#spinner-app-fee").classList.add("hidden");
    document.querySelector("#button-text-app-fee").classList.remove("hidden");
  }
}