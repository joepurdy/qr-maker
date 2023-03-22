<script>
  import { writable } from 'svelte/store';

  let logoFile;
  let qrCodeSize = '300';
  let qrCodeData = '';
  let qrCodeURL = writable(null);

  async function handleSubmit() {
    const formData = new FormData();
    formData.append('logo', logoFile[0]);
    formData.append('size', qrCodeSize);
    formData.append('data', qrCodeData);

    try {
      const response = await fetch('http://localhost:8080/generate-qr-code', {
        method: 'POST',
        body: formData,
      });

      if (response.ok) {
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        $qrCodeURL = url;
      } else {
        console.error('Failed to generate QR code');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  }
</script>

<section class="container mx-auto px-4 mt-10">
  <h1 class="text-2xl font-bold mb-8 text-center">QR Code Generator</h1>
  <form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4" on:submit|preventDefault={handleSubmit}>
    <div class="mb-4">
      <label class="block text-gray-700 text-sm font-bold mb-2">Logo</label>
      <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight" type="file" bind:files={logoFile} accept="image/jpeg,image/png,image/svg+xml" />
    </div>
    <div class="mb-4">
      <label class="block text-gray-700 text-sm font-bold mb-2">QR Code Size</label>
      <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight" type="number" bind:value={qrCodeSize} />
    </div>
    <div class="mb-6">
      <label class="block text-gray-700 text-sm font-bold mb-2">QR Code Data</label>
      <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight" type="text" bind:value={qrCodeData} />
    </div>
    <div class="flex items-center justify-center">
      <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" type="submit">Generate QR Code</button>
    </div>
  </form>
</section>

{#if $qrCodeURL}
  <div class="container mx-auto px-4 mt-8">
    <img class="mx-auto border-2 border-gray-300 rounded" src={$qrCodeURL} alt="Generated QR code" />
  </div>
{/if}
